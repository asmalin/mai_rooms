import React, { useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import CollapsablePanel from "./CollapsablePanel";
import "../styles/Lesson.css";
import checkAuth from "../services/authService.js";
import { reserveRoom, cancelReserve } from "../services/ReservationService.js";

export default function Lesson({ fetchReservedRooms, user_id, ...props }) {
  const [comment, setComment] = useState("");
  const location = useLocation();
  const navigate = useNavigate();

  function reserveRoomHandler(event) {
    event.preventDefault();
    checkAuth().then((user) => {
      if (user === null) {
        navigate("/login", {
          state: { from: location.pathname + location.search },
        });
      }

      const lessonForReserve = {
        roomId: props.roomId,
        date: props.date,
        startTime: props.timeStart,
        endTime: props.timeEnd,
        comment: comment,
      };
      reserveRoom(lessonForReserve).then(() => fetchReservedRooms());
      setComment("");
    });
  }

  function cancelReservationHandler(event) {
    event.preventDefault();
    checkAuth().then((user) => {
      if (user === null) {
        navigate("/login", {
          state: { from: location.pathname + location.search },
        });
      }

      const lessonForCancelReserve = {
        roomId: props.roomId,
        date: props.date,
        startTime: props.timeStart,
      };

      cancelReserve(lessonForCancelReserve).then(() => fetchReservedRooms());
    });
  }

  let panelTitle =
    props.lessonNumber + ". " + props.timeStart + "-" + props.timeEnd + " ";
  if (!props.lesson) {
    return (
      <CollapsablePanel
        title={
          <>
            <span className="lessonTime">{panelTitle}</span>
            <span>Свободно</span>
          </>
        }
        type={"free"}
      >
        <div>
          <textarea
            value={comment}
            placeholder="Комментарий"
            onChange={(event) => setComment(event.target.value)}
            rows="3"
          ></textarea>
          <button onClick={reserveRoomHandler}>Забронировать</button>
        </div>
      </CollapsablePanel>
    );
  } else if (props.lesson.reserver) {
    const { reserver, comment } = props.lesson;
    return (
      <CollapsablePanel
        title={
          <>
            <span className="lessonTime">{panelTitle}</span>
            <span>{reserver}</span>
          </>
        }
        type={"reserved_lesson"}
      >
        <ul>
          <li>
            <strong>Комментарий: </strong>
            {comment || "без объяснения причины!"}
          </li>
        </ul>

        {props.lesson.reserver_id === user_id && (
          <button
            className="cancelReserve-btn"
            onClick={cancelReservationHandler}
          >
            Отменить
          </button>
        )}
      </CollapsablePanel>
    );
  } else {
    const { lector, subject, groups, type } = props.lesson;
    //panelTitle ;

    return (
      <CollapsablePanel
        title={
          <>
            <span className="lessonTime">{panelTitle}</span>
            <span>{lector || "Информация о преподавателе отсутствует"}</span>
          </>
        }
        type={"schedule_lesson"}
      >
        <ul>
          <li>
            {subject} ({type})
          </li>
          <li>{groups}</li>
        </ul>
      </CollapsablePanel>
    );
  }
}
