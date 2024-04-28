import { Link, NavLink, useLocation, useNavigate } from "react-router-dom";
import "../styles/Header.css";

export default function Header({ isLoggedIn, setIsLoggedIn, user }) {
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
          {user?.role === "admin" && (
            <>
              <li className="navbar__item">
                <NavLink
                  to="/reservation_management"
                  className="navbar__link"
                  activeclassname="active"
                >
                  Бронирование
                </NavLink>
              </li>
              <li className="navbar__item">
                <NavLink
                  to="/qrcode"
                  className="navbar__link"
                  activeclassname="active"
                >
                  QRcodes
                </NavLink>
              </li>
            </>
          )}
        </ul>
        {isLoggedIn ? (
          <div className="navbar__auth">
            <Link to="/profile">
              <div className="user">{user?.fullname}</div>
            </Link>

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
    </div>
  );
}
