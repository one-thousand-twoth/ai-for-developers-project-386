import assert from 'node:assert/strict'
import { resolve } from 'node:path'
import { DatabaseSync } from 'node:sqlite'
import { test } from 'node:test'
import { createMeetingType, getHealth, getMeetingTypes, type MeetingType } from '@calendar/sdk'
import { server } from './harness.js'

test('a guest sees a healthy calendar and an empty collection', async t => {
  const app = await server(t)
  const client = await app.start()
  assert.deepEqual((await getHealth({ client })).data, { status: 'ok' })
  const result = await getMeetingTypes({ client })
  assert.equal(result.response?.status, 200)
  assert.deepEqual(result.data, [])
})

test('invalid input is rejected without creating a meeting type', async t => {
  const app = await server(t)
  const client = await app.start()
  const invalid = [
    { name: ' \t\n\u2003\u00a0', durationMinutes: 30 },
    { name: '', durationMinutes: 30 },
    { durationMinutes: 30 },
    { name: 'Call' },
    { name: 'Call', durationMinutes: 0 },
    { name: 'Call', durationMinutes: 481 },
    { name: 'Call', durationMinutes: 1.5 },
    { name: 'Call', durationMinutes: '30' },
    { name: 'Call', durationMinutes: null },
    { name: 'x'.repeat(201), durationMinutes: 30 },
    { name: 'Call', durationMinutes: 30, description: 'x'.repeat(2001) },
    { name: 'Call', durationMinutes: 30, description: null },
    null,
  ]
  for (const body of invalid) {
    const result = await createMeetingType({
      client,
      body: { name: 'placeholder', durationMinutes: 30 },
      bodySerializer: () => JSON.stringify(body),
    })
    assert.equal(result.response?.status, 400, JSON.stringify(body))
    assert.equal(result.error?.code, 'bad_request')
    assert.ok(result.error?.message)
    assert.deepEqual(Object.keys(result.error).sort(), ['code', 'message'])
  }
  for (const body of ['{"name":', '{"name":"Call","durationMinutes":30}{}']) {
    const result = await createMeetingType({
      client, body: { name: 'placeholder', durationMinutes: 30 }, bodySerializer: () => body,
    })
    assert.equal(result.response?.status, 400)
    assert.equal(result.error?.code, 'bad_request')
    assert.ok(result.error?.message)
  }
  assert.deepEqual((await getMeetingTypes({ client })).data, [])
})

test('an owner creates a meeting type and a guest reads it', async t => {
  const app = await server(t)
  const client = await app.start()
  const result = await createMeetingType({ client, body: {
    name: '  Знакомство\u2003', durationMinutes: 30, description: 'Первый звонок',
  } })
  assert.equal(result.response?.status, 201)
  assert.ok(result.data)
  assert.match(result.data.id, /^[0-9a-f-]{36}$/)
  assert.deepEqual(result.data, {
    id: result.data.id, name: 'Знакомство', durationMinutes: 30, description: 'Первый звонок',
  })
  assert.deepEqual((await getMeetingTypes({ client })).data, [result.data])
})

test('names are unique after Unicode case folding and trimming', async t => {
  const app = await server(t)
  const client = await app.start()
  for (const [name, duplicate] of [
    ['Знакомство', '\u2003\tзнаКОМСТВО\u00a0'],
    ['Straße', 'STRASSE'],
    ['Σ', 'ς'],
  ]) {
    const created = await createMeetingType({ client, body: { name, durationMinutes: 30 } })
    assert.equal(created.response?.status, 201)
    const conflict = await createMeetingType({ client, body: {
      name: duplicate, durationMinutes: 60, description: 'Different details do not allow a duplicate name',
    } })
    assert.equal(conflict.response?.status, 409)
    assert.equal(conflict.error?.code, 'duplicate_meeting_type')
    assert.ok(conflict.error?.message)
  }
  assert.equal((await getMeetingTypes({ client })).data?.length, 3)
})

test('concurrent owners on two servers create only one normalized name', async t => {
  const first = await server(t)
  const second = await server(t)
  const clients = [
    await first.start(),
    await second.start({ DATABASE_PATH: resolve(first.directory, 'calendar.sqlite') }),
  ]
  const results = await Promise.all(Array.from({ length: 32 }, (_, index) => createMeetingType({
    client: clients[index % 2],
    body: { name: index % 2 ? '\u2003ВСТРЕЧА\u00a0' : 'встреча', durationMinutes: 30 },
  })))
  assert.equal(results.filter(result => result.response?.status === 201).length, 1)
  const conflicts = results.filter(result => result.response?.status === 409)
  assert.equal(conflicts.length, 31)
  for (const conflict of conflicts) assert.equal(conflict.error?.code, 'duplicate_meeting_type')
  const types = (await getMeetingTypes({ client: clients[0] })).data
  assert.equal(types?.length, 1)
  assert.deepEqual((await getMeetingTypes({ client: clients[1] })).data, types)
  await second.stop()
})

test('boundary values and optional descriptions survive repeated restarts', async t => {
  const app = await server(t)
  let client = await app.start()
  const created: MeetingType[] = []
  for (const body of [
    { name: 'Короткая', durationMinutes: 1 },
    { name: 'Длинная', durationMinutes: 480, description: '' },
    { name: '😀'.repeat(200), durationMinutes: 30, description: '😀'.repeat(2000) },
    { name: 'Два пробела', durationMinutes: 30 },
    { name: 'Два  пробела', durationMinutes: 30 },
  ]) {
    const result = await createMeetingType({ client, body })
    assert.equal(result.response?.status, 201)
    assert.ok(result.data)
    assert.deepEqual(result.data, { ...body, id: result.data.id })
    created.push(result.data)
  }
  const byID = (a: MeetingType, b: MeetingType) => a.id.localeCompare(b.id)
  for (let restart = 0; restart < 2; restart++) {
    await app.stop()
    client = await app.start({ CALENDAR_TIMEZONE: 'Pacific/Auckland' })
    const list = await getMeetingTypes({ client })
    assert.equal(list.response?.status, 200)
    assert.deepEqual(list.data?.sort(byID), [...created].sort(byID))
    const duplicate = await createMeetingType({ client, body: { name: ' КОРОТКАЯ ', durationMinutes: 15 } })
    assert.equal(duplicate.response?.status, 409)
    assert.equal(duplicate.error?.code, 'duplicate_meeting_type')
  }
  const newType = await createMeetingType({ client, body: { name: 'После рестарта', durationMinutes: 60 } })
  assert.equal(newType.response?.status, 201)
  assert.equal((await getMeetingTypes({ client })).data?.length, 6)
})

test('a locked database returns temporary unavailability rather than a name conflict', async t => {
  const app = await server(t)
  const client = await app.start()
  // A real external connection injects a storage outage; assertions stay at the SDK boundary.
  const locker = new DatabaseSync(resolve(app.directory, 'calendar.sqlite'))
  locker.exec('BEGIN EXCLUSIVE')
  try {
    const results = await Promise.all([
      createMeetingType({ client, body: { name: 'Retry later', durationMinutes: 30 } }),
      getMeetingTypes({ client }),
    ])
    for (const result of results) {
      assert.equal(result.response?.status, 503)
      assert.equal(result.error?.code, 'temporarily_unavailable')
      assert.ok(result.error?.message)
    }
  } finally {
    locker.exec('ROLLBACK')
    locker.close()
  }
  assert.deepEqual((await getMeetingTypes({ client })).data, [])
  assert.equal((await createMeetingType({ client, body: { name: 'Retry later', durationMinutes: 30 } })).response?.status, 201)
})
