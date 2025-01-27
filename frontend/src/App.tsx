import './App.css'
import {Route, Routes} from "react-router-dom";
import List from "@/components/pages/List.tsx";
import Interview from "@/components/pages/Interview.tsx";
import {Toaster} from "@/components/ui/toaster.tsx";

function App() {
  return (
      <div className="content">
      <Routes>
          <Route path="*" element={<List />} />
          <Route path="/interview/:interviewId" element={<Interview />}/>
      </Routes>
          <Toaster/>
      </div>
  )
}

export default App
