import { Route, Routes } from "react-router-dom";

import LoginPage from "@/pages/login";
import HomePage from "@/pages/home";

function App() {
  return (
    <div>
      <Routes>
        <Route element={<LoginPage />} path="/" />
        <Route element={<HomePage />} path="/home" />
      </Routes>
    </div>
  );
}

export default App;
