import { Link, NavLink, useLocation, useNavigate } from "react-router-dom";
import "../styles/Header.css";

import { logout } from "../services/authService";

export default function Header({ user, setUser }) {
  const location = useLocation();
  const navigate = useNavigate();

  const loginHandler = () => {
    navigate("/login", {
      state: { from: location.pathname + location.search },
    });
  };

  const logoutHandler = () => {
    localStorage.removeItem("token");
    logout();
    setUser(null);
  };

  return (
    <div className="wrapper">
      <nav className="navbar">
        <ul className="navbar__list">
          <li className="navbar__item">
            <NavLink to="/" className="navbar__link" activeclassname="active">
              MAI rooms
            </NavLink>
          </li>
          <li className="navbar__item">
            <NavLink
              to="/info"
              className="navbar__link"
              activeclassname="active"
            >
              Инфо
            </NavLink>
          </li>
        </ul>
        {user ? (
          <div className="navbar__auth">
            <Link to="/profile">
              <div className="user">{user?.username}</div>
            </Link>

            <button className="logout-btn" onClick={logoutHandler}>
              Выйти
            </button>
          </div>
        ) : (
          location.pathname !== "/login" && (
            <div className="navbar__auth">
              <button className="login-btn" onClick={loginHandler}>
                Войти
              </button>
            </div>
          )
        )}
      </nav>
    </div>
  );
}
