import React from "react";
import "../styles/Footer.css";
import tgLogo from "../img/tg_logo.png";

export default function Footer() {
  return (
    <>
      <footer>
        <div className="year">{new Date().getFullYear()}</div>
        <div className="social">
          <a
            className="btn btn-gray btn-round"
            href="https://t.me/roomsReservation_bot"
          >
            <img className="logo" src={tgLogo} alt="Telegram Logo" />
          </a>
        </div>
      </footer>
    </>
  );
}
