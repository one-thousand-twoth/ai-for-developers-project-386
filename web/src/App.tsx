import { useState } from 'react'
import './App.css'

function App() {
  const [count, setCount] = useState(0)

  return (
    <main className="welcome">
      <p className="eyebrow">Календарь звонков</p>
      <h1>Время для встречи</h1>
      <p className="description">
        Здесь начинается приложение для планирования и бронирования звонков.
      </p>
      <button type="button" onClick={() => setCount((value) => value + 1)}>
        Проверить React: {count}
      </button>
      <p className="hint" role="status">Нажатий: {count}</p>
    </main>
  )
}

export default App
