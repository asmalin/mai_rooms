import React, { useState } from "react";
import { useLocation } from "react-router-dom";
import Room from "../components/Room.js";
import "../styles/Rooms.css";

export default function Rooms({ user_id }) {
  const [showOnlyFree, setShowOnlyFree] = useState(false);

  const searchParams = new URLSearchParams(useLocation().search);
  const roomIds = searchParams.getAll("roomId");
  const date = new Date(searchParams.get("date")).toLocaleDateString("ru");

  if (roomIds.length === 0) return <div>Не выбрана аудитория!</div>;

  if (date === "Invalid Date") return <div>Неверный формат даты!</div>;

  return (
    <>
      <div className="settingChanges">
        <div>{date}</div>
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
            roomId={id}
            date={date}
            showOnlyFree={showOnlyFree}
            user_id={user_id}
          />
        </div>
      ))}
    </>
  );
}
