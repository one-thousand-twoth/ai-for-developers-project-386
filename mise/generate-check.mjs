import { readdir, readFile } from 'node:fs/promises'
import { spawnSync } from 'node:child_process'

const roots = ['openapi', 'gen/api', 'web/src/generated/api']
async function snapshot() {
  const files = new Map()
  async function visit(path) {
    for (const entry of await readdir(path, { withFileTypes: true }).catch((error) => {
      if (error.code === 'ENOENT') return []
      throw error
    })) {
      const name = `${path}/${entry.name}`
      if (entry.isDirectory()) await visit(name)
      else files.set(name, await readFile(name, 'base64'))
    }
  }
  for (const root of roots) await visit(root)
  return JSON.stringify([...files].sort(([a], [b]) => a.localeCompare(b)))
}
const before = await snapshot()
const generated = spawnSync('mise', ['run', 'generate'], { stdio: 'inherit' })
if (generated.status !== 0) process.exit(generated.status ?? 1)
const changed = before !== await snapshot()
const diff = spawnSync('git', ['diff', '--exit-code', '--', ...roots], { stdio: 'inherit' })
const untracked = spawnSync('git', ['ls-files', '--others', '--exclude-standard', '--', ...roots], { encoding: 'utf8' })
if (diff.status !== 0 || untracked.status !== 0 || changed || untracked.stdout.trim()) {
  console.error('Generated artifacts differ or are untracked. Run mise run generate and commit all artifacts.')
  process.exit(1)
}
