import { cancelReserve } from "../services/ReservationService.js";
import React, { useEffect, useState } from "react";
import axios from "axios";
import AdminNav from "../components/AdminNav.js";
import "../styles/ReservationMgmt.css";

export default function ReservationMgmt({ user }) {
  const [reservedLessons, setReservedLessons] = useState();
  const [filter, setFilter] = useState({
    room_name: "",
    reserver: "",
    date: "",
  });

  const [showConfirm, setShowConfirm] = useState(false);

  const [lessonForCanceling, setLessonForCanceling] = useState({
    roomId: 0,
    roomName: "",
    date: "",
    startTime: "",
    endTime: "",
  });

  const confirmCancelLesson = (roomId, roomName, date, timeStart, timeEnd) => {
    setShowConfirm(true);
    setLessonForCanceling({
      roomId: roomId,
      roomName: roomName,
      date: date,
      startTime: timeStart,
      timeEnd: timeEnd,
    });
  };

  useEffect(() => {
    fetchReservedRooms();
  }, []);

  function fetchReservedRooms() {
    axios
      .get("http://localhost:5001/api/all_reserved_lesssons")
      .then((response) => {
        setReservedLessons(response.data);
      })
      .catch((error) => {
        console.error("There was an error!", error);
      });
  }

  const filteredLessons = reservedLessons
    ? reservedLessons.filter((lesson) => {
        return (
          lesson.room_name
            .toLowerCase()
            .includes(filter.room_name.toLowerCase()) &&
          lesson.reserver
            .toLowerCase()
            .includes(filter.reserver.toLowerCase()) &&
          lesson.date.includes(formatDate(filter.date))
        );
      })
    : [];

  function handleCancelLesson(roomId, date, time_start) {
    setShowConfirm(false);
    const lessonForCancelReserve = {
      roomId: roomId,
      date: date,
      startTime: time_start,
    };

    cancelReserve(lessonForCancelReserve).then(() => fetchReservedRooms());
  }

  const handleOverlayClick = (e) => {
    if (e.target === e.currentTarget) {
      setShowConfirm(false);
    }
  };

  return (
    <>
      <AdminNav />
      <div className="reservation_mgmt ">
        <strong>Все заронированные аудитории</strong>

        <div>
          <strong>Фильтры</strong>
          <form className="filter_form">
            <input
              type="text"
              placeholder="Аудитория"
              value={filter.room_name}
              onChange={(e) =>
                setFilter({ ...filter, room_name: e.target.value })
              }
            />
            <input
              type="text"
              placeholder="Имя"
              value={filter.reserver}
              onChange={(e) =>
                setFilter({ ...filter, reserver: e.target.value })
              }
            />
            <input
              className="date_input"
              type="date"
              placeholder="Дата"
              value={filter.date}
              onChange={(e) => setFilter({ ...filter, date: e.target.value })}
            />
          </form>
        </div>
        {filteredLessons.length > 0 ? (
          <table className="table table-bordered">
            <thead>
              <tr>
                <th scope="col">#</th>
                <th scope="col">classroom</th>
                <th scope="col">reserver name</th>
                <th scope="col">date</th>
                <th scope="col">time start</th>
                <th scope="col">time end</th>
                <th scope="col">comment</th>
                <th scope="col"></th>
              </tr>
            </thead>
            <tbody>
              {filteredLessons.map((lesson, index) => (
                <tr key={index}>
                  <th scope="row">{index + 1}</th>
                  <td>{lesson.room_name}</td>
                  <td>{lesson.reserver}</td>
                  <td>{lesson.date}</td>
                  <td>{lesson.time_start}</td>
                  <td>{lesson.time_end}</td>
                  <td>{lesson.comment}</td>
                  <td align="center">
                    <button
                      className="btn btn-danger"
                      onClick={() =>
                        confirmCancelLesson(
                          lesson.room_id,
                          lesson.room_name,
                          lesson.date,
                          lesson.time_start,
                          lesson.time_end
                        )
                      }
                    >
                      Отменить
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : (
          <h2>Список пуст!</h2>
        )}
        {showConfirm && (
          <div className="overlay" onClick={handleOverlayClick}>
            <div className="confirmation-dialog">
              <p>
                Вы уверены, что хотите удалить бронь аудитории{" "}
                <b>{lessonForCanceling.roomName}</b> {lessonForCanceling.date} с{" "}
                {lessonForCanceling.startTime} по {lessonForCanceling.timeEnd}?
              </p>

              <div className="confirmDlgBtns">
                <button
                  className="btn btn-danger"
                  onClick={() =>
                    handleCancelLesson(
                      lessonForCanceling.roomId,
                      lessonForCanceling.date,
                      lessonForCanceling.startTime
                    )
                  }
                >
                  Да
                </button>
                <button
                  className="btn btn-secondary"
                  onClick={() => setShowConfirm(false)}
                >
                  Нет
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </>
  );
}

function formatDate(time) {
  if (time === "") return "";
  const [year, month, day] = time.split("-");
  return `${day}.${month}.${year}`;
}
