import './App.css'
import {Route, Routes} from "react-router-dom";
import List from "@/components/pages/List.tsx";
import Interview from "@/components/pages/Interview.tsx";

function App() {
  return (
      <Routes>
          <Route path="*" element={<List />} />
          <Route path="/interview/:interviewId" element={<Interview />}/>
      </Routes>
  )
}

export default App
