import React from "react";
import "../styles/qrcode.css";
import AdminNav from "../components/AdminNav";
import qrcodeLogo from "../img/qrcode_scan.png";
import { CreateQRCodes } from "../services/qrcodeService";
import { QRCodeGenerator } from "../services/qrcodeService";
import { useEffect, useState } from "react";
import axios from "axios";

export default function QRcodesGenerator({ user }) {
  const [buildingsList, setBuildingsList] = useState([]);
  const [selectedBuilding, setSelectedBuilding] = useState();
  const [roomsList, setRoomsList] = useState([]);
  const [selectedRooms, setSelectedRooms] = useState([]);
  const [error, setError] = useState(false);

  const [showClickedQRcode, setShowClickedQRcode] = useState(false);
  const [roomIdForQRcode, setroomIdForQRcode] = useState(false);

  const handleOverlayClick = (e) => {
    if (e.target === e.currentTarget) {
      setShowClickedQRcode(false);
    }
  };

  const showQRcode = (roomId) => {
    setroomIdForQRcode(roomId);
    setShowClickedQRcode(true);
  };

  useEffect(() => {
    axios
      .get("http://localhost:5001/api/buildings")
      .then((response) => {
        setBuildingsList(response.data);
      })
      .catch((error) => {
        console.error("There was an error!", error);
      });
  }, []);

  useEffect(() => {
    if (selectedBuilding) {
      const apiUrl = "http://localhost:5001/api/rooms/" + selectedBuilding;
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

    CreateQRCodes(selectedRooms);
  };

  return (
    <>
      <AdminNav />
      <div className="qrcode_generator">
        <h2>Создать QRCode для выбранной аудиториий</h2>

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
                  <img
                    className="qrcode_logo"
                    src={qrcodeLogo}
                    alt="QRcode scan logo"
                    onClick={() => showQRcode(room.id)}
                  />
                </div>
              ))}
            </div>
            {error && (
              <div className="emptySelectedRoomsError">
                Выберете хотя бы одну аудиторию
              </div>
            )}
            <button type="submit" className="selectRoomButton">
              Скачать
            </button>
          </form>
        )}
        {showClickedQRcode && (
          <div className="overlay" onClick={handleOverlayClick}>
            <QRCodeGenerator roomId={roomIdForQRcode} />
          </div>
        )}
      </div>
    </>
  );
}
