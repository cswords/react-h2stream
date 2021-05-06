import './App.css';
import { useEffect, useState } from 'react'

const go = new window.global.Go();

function App() {
  const [wasmResourse, setWasmResourse] = useState(null)
  const [sum, setSum] = useState(0)
  const [count, setCount] = useState(0)

  useEffect(() => {
    WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then(result => {
      if (!wasmResourse && result.instance) {
        setWasmResourse(result)
        go.run(result.instance)
      }
    })
  }, [wasmResourse])

  return (
    <div className="App">
      {wasmResourse && <button onClick={(e) => {
        window.global.getH2C('https://localhost:8443', (data) => {
          console.log('[JS] Got data: ', data)
          const f = data.substring(6)
          setSum(sum + parseFloat(f))
          setCount(count + 1)
        })
      }}>Call</button>}
      <div>
        {count === 0 ? sum : sum / count}
      </div>
    </div>
  );
}

export default App;
