import { useState, type FormEvent } from 'react'
import { createBooking, getBooking, type Booking } from './generated/api'
import './App.css'

function errorMessage(error: unknown): string {
  if (error && typeof error === 'object' && 'message' in error) return String(error.message)
  return 'Не удалось связаться с сервером. Проверьте, запущен ли backend.'
}

function App() {
  const [booking, setBooking] = useState<Booking | null>(null)
  const [id, setId] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  async function submit(event: FormEvent<HTMLFormElement>, action: 'create' | 'find') {
    event.preventDefault()
    const form = new FormData(event.currentTarget)
    setLoading(true)
    setError('')
    setBooking(null)
    try {
      const response = action === 'create'
        ? await createBooking({ body: {
          title: String(form.get('title')),
          startsAt: new Date(String(form.get('startsAt'))).toISOString(),
          endsAt: new Date(String(form.get('endsAt'))).toISOString(),
        } })
        : await getBooking({ path: { id: id.trim() } })
      if (response.error) {
        setError(errorMessage(response.error))
      } else if (response.data) {
        setBooking(response.data)
        setId(response.data.id)
      } else {
        setError('Сервер вернул пустой ответ.')
      }
    } catch (caught) {
      setError(errorMessage(caught))
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="welcome">
      <p className="eyebrow">Календарь звонков · демонстрация API</p>
      <h1>Время для встречи</h1>
      <p className="description">Создайте бронирование или найдите его по ID. Данные хранятся до перезапуска сервера.</p>
      <form onSubmit={(event) => void submit(event, 'create')}>
        <h2>Создать бронирование</h2>
        <label>Название<input name="title" required pattern=".*\S.*" /></label>
        <label>Начало<input name="startsAt" type="datetime-local" required /></label>
        <label>Конец<input name="endsAt" type="datetime-local" required /></label>
        <p className="hint">Время указано в часовом поясе вашего устройства.</p>
        <button disabled={loading}>Создать</button>
      </form>
      <form onSubmit={(event) => void submit(event, 'find')}>
        <h2>Найти по ID</h2>
        <label>ID бронирования<input name="id" value={id} onChange={(event) => setId(event.target.value)} required pattern=".*\S.*" /></label>
        <button disabled={loading}>Найти</button>
      </form>
      <div role="status" aria-live="polite">{loading && 'Загрузка…'}</div>
      {error && <p role="alert" className="error">{error}</p>}
      {booking && <section aria-label="Результат" aria-live="polite">
        <h2>{booking.title}</h2>
        <p className="booking-id">ID: {booking.id}</p>
        <p>Начало: <time dateTime={booking.startsAt}>{new Date(booking.startsAt).toLocaleString()}</time></p>
        <p>Конец: <time dateTime={booking.endsAt}>{new Date(booking.endsAt).toLocaleString()}</time></p>
      </section>}
    </main>
  )
}

export default App
