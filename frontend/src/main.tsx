import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import {BrowserRouter} from "react-router-dom";
import {Network} from "@/lib/network/network.ts";
import './code-highlight.css';


let addr = 'http://localhost:9595';
let isDev = false;

if (process.env.NODE_ENV === 'development') {
    isDev = true;
    addr = "http://localhost:9595";

    console.log("working in dev mode");
}

const network = new Network(addr, isDev);
globalThis.network = network;

createRoot(document.getElementById('root')!).render(
  <StrictMode>
      <BrowserRouter>
            <App />
      </BrowserRouter>
  </StrictMode>,
)
