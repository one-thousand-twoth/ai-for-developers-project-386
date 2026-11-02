import assert from 'node:assert/strict'
import { spawn, type ChildProcess } from 'node:child_process'
import { once } from 'node:events'
import { mkdtemp, rm } from 'node:fs/promises'
import { resolve } from 'node:path'
import { setTimeout as delay } from 'node:timers/promises'
import type { TestContext } from 'node:test'
import { getHealth } from '@calendar/sdk'
import { createClient } from '@calendar/sdk/client'

export async function server(t: TestContext) {
  const directory = await mkdtemp(resolve('.test/api/db-'))
  let child: ChildProcess | undefined
  let logs = ''

  async function stop() {
    if (!child?.pid || child.exitCode !== null || child.signalCode !== null) return
    const exited = once(child, 'exit')
    child.kill('SIGTERM')
    const timeout = setTimeout(() => child?.kill('SIGKILL'), 3_000)
    try { await exited } finally { clearTimeout(timeout) }
  }

  t.after(async () => {
    try {
      await stop()
    } finally {
      t.diagnostic(logs)
      await rm(directory, { recursive: true, force: true })
    }
  })

  async function start(overrides: Record<string, string> = {}) {
    assert.ok(!child || child.exitCode !== null || child.signalCode !== null, 'stop before restart')
    logs = ''
    child = spawn(resolve('.test/api/server'), [], {
      cwd: directory,
      env: {
        ...process.env,
        DATABASE_PATH: resolve(directory, 'calendar.sqlite'),
        CALENDAR_TIMEZONE: 'Europe/Moscow',
        LISTEN_ADDR: '127.0.0.1:0',
        ...overrides,
      },
      stdio: ['ignore', 'pipe', 'pipe'],
    })
    let spawnError: Error | undefined
    child.on('error', error => { spawnError = error })
    child.stdout?.on('data', chunk => { logs += chunk.toString() })
    child.stderr?.on('data', chunk => { logs += chunk.toString() })
    const deadline = Date.now() + 10_000
    while (Date.now() < deadline) {
      if (spawnError) throw spawnError
      if (child.exitCode !== null || child.signalCode !== null) throw new Error(`server exited:\n${logs}`)
      const address = /server listening on (http:\/\/127\.0\.0\.1:\d+)/.exec(logs)?.[1]
      if (address) {
        const client = createClient({ baseUrl: address, signal: t.signal })
        try {
          const health = await getHealth({ client, signal: AbortSignal.timeout(500) })
          if (health.data?.status === 'ok') return client
        } catch { /* Retry until the bounded startup deadline. */ }
      }
      await delay(20)
    }
    throw new Error(`health timeout:\n${logs}`)
  }

  return { start, stop, directory }
}
