import "./App.css";
import Header from "./components/Header.js";
import Footer from "./components/Footer.js";
import Home from "./pages/home.js";
import Info from "./pages/info.js";
import Rooms from "./pages/rooms.js";
import Login from "./pages/login.js";
import Profile from "./pages/profile.js";
import Users from "./pages/Users.js";
import ReservationMgmt from "./pages/reservationMgmt.js";
import QRcodesGenerator from "./pages/qrcodeGenerator.js";

import Error403 from "./pages/errorPages/error403.js";

import { Route, Routes } from "react-router-dom";
import { useEffect, useState } from "react";

import { checkAuth, refreshTokens } from "./services/authService.js";

function App() {
  const [user, setUser] = useState();

  useEffect(() => {
    const fetchData = async () => {
      if (!localStorage.getItem("token")) return;

      let response = await tryToGetUser();
      if (response === 401) {
        response = await refreshTokens();
        if (response.ok) {
          const data = await response.json();
          localStorage.setItem("token", JSON.stringify(data.token));
          await tryToGetUser();
        }
      }
    };

    fetchData();
  }, []);

  const tryToGetUser = async () => {
    const response = await checkAuth();
    if (response.ok) {
      setUser(await response.json());
    }
    return response.status;
  };

  return (
    <>
      <Header user={user} setUser={setUser} />
      <main>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/rooms" element={<Rooms user={user} />} />
          <Route path="/info" element={<Info />} />
          <Route path="/login" element={<Login setUser={setUser} />} />
          <Route
            path="/profile"
            element={
              user ? <Profile user={user} /> : <Login setUser={setUser} />
            }
          />
          <Route
            path="/profile/reservation_management"
            element={
              user?.role === "admin" ? (
                <ReservationMgmt user={user} />
              ) : (
                <Error403 />
              )
            }
          />
          <Route
            path="/profile/qrcode"
            element={
              user?.role === "admin" ? (
                <QRcodesGenerator user={user} />
              ) : (
                <Error403 />
              )
            }
          />
          <Route
            path="/profile/users"
            element={
              user?.role === "admin" ? <Users user={user} /> : <Error403 />
            }
          />
        </Routes>
      </main>
      <Footer />
    </>
  );
}

export default App;
