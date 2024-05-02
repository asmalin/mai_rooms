import "../styles/home.css";
import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";

function Home() {
  const [buildingsList, setBuildingsList] = useState([]);
  const [selectedBuilding, setSelectedBuilding] = useState();
  const [roomsList, setRoomsList] = useState([]);
  const [selectedRooms, setSelectedRooms] = useState([]);
  const [error, setError] = useState(false);
  const [selectedDate, setSelectedDate] = useState(
    new Date()
      .toLocaleDateString("ru-RU", { timeZone: "Europe/Moscow" })
      .split(".")
      .reverse()
      .join("-")
  );

  const navigate = useNavigate();

  const handleDateChange = (event) => {
    setSelectedDate(event.target.value);
  };

  useEffect(() => {
    const server_domain = process.env.REACT_APP_SERVER_BASE_URL;
    axios
      .get(server_domain + "/api/buildings")
      .then((response) => {
        setBuildingsList(response.data);
      })
      .catch((error) => {
        console.error("There was an error!", error);
      });
  }, []);

  useEffect(() => {
    if (selectedBuilding) {
      const server_domain = process.env.REACT_APP_SERVER_BASE_URL;
      const apiUrl = server_domain + "/api/rooms/" + selectedBuilding;
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

  const getRooms = (event) => {
    event.preventDefault();

    if (selectedRooms.length === 0) {
      setError(true);
      return;
    }

    const formData = new FormData(event.target);
    const queryString = new URLSearchParams(formData).toString();
    const url = `/rooms?${queryString}`;
    navigate(url);
  };

  return (
    <div className="wrapper">
      <h1>Найти и забронировать свободную аудиторию МАИ</h1>
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
        <form onSubmit={getRooms} className="roomSelectForm">
          <input
            className="date_input"
            type="date"
            pattern="\d{4}-\d{2}-\d{2}"
            name="date"
            value={selectedDate}
            onChange={handleDateChange}
            required
          />
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
            <div className="errorMsg">Выберете хотя бы одну аудиторию</div>
          )}
          <button type="submit" className="selectRoomButton">
            Перейти
          </button>
        </form>
      )}
    </div>
  );
}

export default Home;
