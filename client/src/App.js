import { BrowserRouter, Route, Router, Routes } from "react-router-dom";
import "./App.css";
import SignUpPage from "./Components/SignUpPage/SignUpPage";
import DashboardPage from "./Components/DashboardPage/DashboardPage";

function App() {
  return (
    <Routes>
      {/* <Route path="/" component={<SignUpPage />} /> */}
      <Route path="/" element={<DashboardPage />} />
    </Routes>
  );
}

export default App;
