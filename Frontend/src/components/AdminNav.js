import React from "react";
import { NavLink } from "react-router-dom";

export default function AdminNav() {
  return (
    <div className="admin_panel">
      <h2>Администрирование</h2>
      <ul>
        <li>
          <NavLink
            to="/profile/reservation_management"
            className="navbar__link"
            activeclassname="active"
          >
            Бронирование
          </NavLink>
        </li>
        <li>
          <NavLink
            to="/profile/qrcode"
            className="navbar__link"
            activeclassname="active"
          >
            QRcodes
          </NavLink>
        </li>
        <li>
          <NavLink
            to="/profile/users"
            className="navbar__link"
            activeclassname="active"
          >
            Пользователи
          </NavLink>
        </li>
      </ul>
    </div>
  );
}
