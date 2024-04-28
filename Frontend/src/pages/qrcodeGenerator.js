import React from "react";
import "../styles/home.css";

import { CreateQRCodes } from "../services/ReservationService";

import { useEffect, useState } from "react";
import axios from "axios";

export default function QRcodesGenerator() {
  const [buildingsList, setBuildingsList] = useState([]);
  const [selectedBuilding, setSelectedBuilding] = useState();
  const [roomsList, setRoomsList] = useState([]);
  const [selectedRooms, setSelectedRooms] = useState([]);
  const [error, setError] = useState(false);

  useEffect(() => {
    axios
      .get("http://localhost:8080/api/buildings")
      .then((response) => {
        setBuildingsList(response.data);
      })
      .catch((error) => {
        console.error("There was an error!", error);
      });
  }, []);

  useEffect(() => {
    if (selectedBuilding) {
      const apiUrl = "http://localhost:8080/api/rooms/" + selectedBuilding;
      fetch(apiUrl)
        .then((response) => response.json())
        .then((data) => setRoomsList(data));
    }
  }, [selectedBuilding]);

  const roomChange = (event) => {
    const roomId = event.target.value;
    const isChecked = event.target.checked;

    if (isChecked) {
      setSelectedRooms([...selectedRooms, roomId]);
    } else {
      setSelectedRooms(selectedRooms.filter((room) => room !== roomId));
    }
    setError(false);
  };

  const createQRcodes = (event) => {
    event.preventDefault();

    CreateQRCodes(selectedRooms).then((res) => console.log(res));
  };

  return (
    <div className="wrapper">
      <h2>Создать QRCode выбранной аудитории для быстрого доступа к ней</h2>
      <div className="select-building">
        <h2>Выберете корпус</h2>
        <select
          className="form-select"
          onChange={(e) => setSelectedBuilding(e.target.value)}
          value={selectedBuilding}
        >
          <option>--Выберете корпус--</option>
          {buildingsList.map((building) => (
            <option key={building.id} value={building.id}>
              {building.name}
            </option>
          ))}
        </select>
      </div>

      {roomsList.length > 0 && (
        <form onSubmit={createQRcodes} className="roomSelectForm">
          <div className="roomsList form-check">
            {roomsList.map((room) => (
              <div className="btn-room" key={room.id}>
                <input
                  id={room.id}
                  type="checkbox"
                  name="roomId"
                  value={room.id}
                  onChange={roomChange}
                />
                <label className="form-check-label" htmlFor={room.id}>
                  {room.name}
                </label>
              </div>
            ))}
          </div>
          {error && (
            <div className="emptySelectedRoomsError">
              Выберете хотя бы одну аудиторию
            </div>
          )}
          <button type="submit" className="selectRoomButton">
            Создать
          </button>
        </form>
      )}
    </div>
  );
}
