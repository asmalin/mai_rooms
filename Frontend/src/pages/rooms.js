import React, { useState } from "react";
import { useLocation } from "react-router-dom";
import Room from "../components/Room.js";
import "../styles/Rooms.css";

export default function Rooms({ user }) {
  const [showOnlyFree, setShowOnlyFree] = useState(false);

  const searchParams = new URLSearchParams(useLocation().search);
  const roomIds = searchParams.getAll("roomId");
  const dateFromUrl = searchParams.get("date");

  const date =
    dateFromUrl != null ? new Date(dateFromUrl).toLocaleDateString("ru") : null;

  if (roomIds.length === 0) return <div>Не выбрана аудитория!</div>;

  if (date === "Invalid Date") return <div>Неверный формат даты!</div>;

  return (
    <>
      <div className="settingChanges">
        <div>{date || getCurrentDate()}</div>
        <div>
          <input
            type="checkbox"
            id="onlyFree"
            onClick={() => {
              setShowOnlyFree(!showOnlyFree);
            }}
          ></input>
          <label htmlFor="onlyFree">Показать только сободные</label>
        </div>
      </div>
      {roomIds.map((id) => (
        <div key={id} className="room">
          <Room
            roomId={Number(id)}
            date={date}
            showOnlyFree={showOnlyFree}
            user={user}
          />
        </div>
      ))}
    </>
  );
}

function getCurrentDate() {
  const date = new Date();
  const day = String(date.getDate()).padStart(2, "0");
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const year = date.getFullYear();
  return `${day}.${month}.${year}`;
}
