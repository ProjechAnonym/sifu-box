import { Route, Routes } from "react-router-dom";

import LoginPage from "@/pages/login";


function App() {
  return (
    <div>
      <Routes>
        <Route element={<LoginPage />} path="/" />
        <Route element={<>555</>} path="/home" />
      </Routes>
    </div>
  );
}

export default App;
