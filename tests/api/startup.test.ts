import assert from 'node:assert/strict'
import { writeFile } from 'node:fs/promises'
import { resolve } from 'node:path'
import { DatabaseSync } from 'node:sqlite'
import { test } from 'node:test'
import { createMeetingType, getMeetingTypes } from '@calendar/sdk'
import { server } from './harness.js'

test('startup requires a named IANA calendar time zone', async t => {
  for (const timeZone of ['', 'Local', '+03:00', 'Mars/Olympus', ' Europe/Moscow ']) {
    await t.test(JSON.stringify(timeZone), async t => {
      const app = await server(t)
      await assert.rejects(app.start({ CALENDAR_TIMEZONE: timeZone }), /CALENDAR_TIMEZONE.*named IANA/)
    })
  }
})

test('the default database path persists meeting types and accepts UTC', async t => {
  const app = await server(t)
  const client = await app.start({ DATABASE_PATH: '', CALENDAR_TIMEZONE: 'UTC' })
  const created = await createMeetingType({ client, body: { name: 'Default database', durationMinutes: 30 } })
  assert.equal(created.response?.status, 201)
  await app.stop()
  const restarted = await app.start({ DATABASE_PATH: resolve(app.directory, 'data/calendar.sqlite'), CALENDAR_TIMEZONE: 'Etc/UTC' })
  assert.deepEqual((await getMeetingTypes({ client: restarted })).data, [created.data])
})

test('startup fails if the database path is a directory', async t => {
  const app = await server(t)
  await assert.rejects(app.start({ DATABASE_PATH: app.directory }), /open database/)
})

test('startup fails for a corrupt database', async t => {
  const app = await server(t)
  await writeFile(resolve(app.directory, 'calendar.sqlite'), 'not a SQLite database')
  await assert.rejects(app.start(), /file is not a database/)
})

test('a failed migration prevents the server from listening', async t => {
  const app = await server(t)
  const database = new DatabaseSync(resolve(app.directory, 'calendar.sqlite'))
  database.exec('CREATE TABLE meeting_types (incompatible TEXT)')
  database.close()
  await assert.rejects(app.start(), /migrate database:.*already exists/)
})
