import "./App.css";
import Header from "./components/Header.js";
import Footer from "./components/Footer.js";
import Home from "./pages/home.js";
import Info from "./pages/info.js";
import Rooms from "./pages/rooms.js";
import Login from "./pages/login.js";
import Profile from "./pages/profile.js";
import ReservationMgmt from "./pages/reservationMgmt.js";
import QRcodesGenerator from "./pages/qrcodeGenerator.js";

import Error403 from "./pages/errorPages/error403.js";

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
        user={user}
      />
      <main>
        <div className="wrapper">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route
              path="/rooms"
              element={<Rooms user={user ? user : null} />}
            />
            <Route path="/info" element={<Info />} />
            <Route
              path="/login"
              element={<Login setIsLoggedIn={setIsLoggedIn} />}
            />
            <Route path="/profile" element={<Profile user={user} />} />
            <Route
              path="/reservation_management"
              element={
                user?.role === "admin" ? (
                  <ReservationMgmt user={user} />
                ) : (
                  <Error403 />
                )
              }
            />
            <Route
              path="/qrcode"
              element={
                user?.role === "admin" ? <QRcodesGenerator /> : <Error403 />
              }
            />
          </Routes>
        </div>
      </main>
      <Footer />
    </>
  );
}

export default App;
