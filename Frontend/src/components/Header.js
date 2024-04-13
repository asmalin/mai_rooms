import { Link, useLocation, useNavigate } from "react-router-dom";
import "../styles/Header.css";

export default function Header({ isLoggedIn, setIsLoggedIn, userFullName }) {
  const location = useLocation();
  const navigate = useNavigate();

  const login = () => {
    navigate("/login", {
      state: { from: location.pathname + location.search },
    });
  };

  const logout = () => {
    localStorage.removeItem("token");
    setIsLoggedIn(false);
  };

  return (
    <>
      <nav className="navbar wrapper">
        <ul className="navbar__list">
          <li className="navbar__item">
            <Link to="/" className="navbar__link">
              MAI rooms
            </Link>
          </li>
          <li className="navbar__item">
            <Link to="/info" className="navbar__link">
              Инфо
            </Link>
          </li>
        </ul>
        {isLoggedIn ? (
          <div className="navbar__auth">
            <div className="user">{userFullName}</div>
            <button className="logout-btn" onClick={logout}>
              Выйти
            </button>
          </div>
        ) : (
          location.pathname !== "/login" && (
            <div className="navbar__auth">
              <button className="login-btn" onClick={login}>
                Войти
              </button>
            </div>
          )
        )}
      </nav>
    </>
  );
}
