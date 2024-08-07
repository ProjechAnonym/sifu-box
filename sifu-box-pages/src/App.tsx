import { Route, Routes } from "react-router-dom";
import DefaultLayout from "@/layout/defaultLayout";
import Login from "@/pages/login";
import Home from "@/pages/home";
import Proxy from "@/pages/proxy";
import Setting from "@/pages/setting";
import "bootstrap-icons/font/bootstrap-icons.css";
function App() {
  return (
    <DefaultLayout>
      <Routes>
        <Route element={<Home />} path="/"></Route>
        <Route element={<Login />} path="/login"></Route>
        <Route element={<Proxy />} path="/proxy"></Route>
        <Route element={<Setting />} path="/setting"></Route>
      </Routes>
    </DefaultLayout>
  );
}

export default App;
