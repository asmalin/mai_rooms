import React, { useState, useEffect } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import "../styles/Login.css";
import axios from "axios";

const API_URL = "http://localhost:8080";

const Login = ({ setIsLoggedIn }) => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [loginErrorMsg, setLoginErrorMsg] = useState("");

  let navigate = useNavigate();
  const location = useLocation();
  const fromPage = location.state?.from;

  useEffect(() => {
    setLoginErrorMsg("");
  }, [username, password]);

  const handleSignIn = async (e) => {
    e.preventDefault();

    await axios
      .post(`${API_URL}/login`, {
        username,
        password,
      })
      .then((response) => {
        setLoginErrorMsg("");

        if (response.data.token) {
          localStorage.setItem("token", JSON.stringify(response.data.token));
          setIsLoggedIn(true);
          if (fromPage) navigate(fromPage);
          else navigate("/");
        }
      })
      .catch((error) => {
        console.log(error);
        setLoginErrorMsg("Неправильный имя пользователя или пароль!");
      });
  };

  return (
    <>
      <form onSubmit={handleSignIn} className="login-form">
        <h2>Авторизация</h2>
        <input
          type="text"
          id="username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          placeholder="username"
        />

        <input
          type="password"
          id="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder="password"
        />
        {loginErrorMsg && <div className="errMessage">{loginErrorMsg}</div>}
        <button type="submit">Войти</button>
      </form>
    </>
  );
};

export default Login;
