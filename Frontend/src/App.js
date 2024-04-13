import "./App.css";
import Header from "./components/Header.js";
import Footer from "./components/Footer.js";
import Home from "./pages/home.js";
import Info from "./pages/info.js";
import Rooms from "./pages/rooms.js";
import Login from "./pages/login.js";
import { Route, Routes } from "react-router-dom";
import { useEffect, useState } from "react";
import checkAuth from "../src/services/authService.js";

function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [user, setUser] = useState();

  useEffect(() => {
    if (localStorage.getItem("token")) {
      checkAuth().then((userData) => {
        if (userData === null) return;
        setUser(userData);
        setIsLoggedIn(true);
      });
    }
  }, [isLoggedIn]);

  return (
    <>
      <Header
        isLoggedIn={isLoggedIn}
        setIsLoggedIn={setIsLoggedIn}
        userFullName={user?.fullname}
      />
      <main>
        <div className="wrapper">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/rooms" element={<Rooms />} />
            <Route path="/info" element={<Info />} />
            <Route
              path="/login"
              element={<Login setIsLoggedIn={setIsLoggedIn} />}
            />
          </Routes>
        </div>
      </main>
      <Footer />
    </>
  );
}

export default App;
