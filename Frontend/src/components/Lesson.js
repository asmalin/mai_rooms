import React, { useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import CollapsablePanel from "./CollapsablePanel";
import "../styles/Lesson.css";
import checkAuth from "../services/authService.js";
import reserveRoom from "../services/ReservationService.js";

export default function Lesson(props) {
  const [comment, setComment] = useState("");
  const location = useLocation();
  const navigate = useNavigate();

  if (!props.lesson) return;

  function submitHandler(event) {
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
        startTime: props.lesson.time_start,
        comment: comment,
      };
      reserveRoom(lessonForReserve).then((response) => console.log(response));
    });
  }

  let panelTitle =
    props.lessonNumber +
    ". " +
    props.lesson.time_start +
    "-" +
    props.lesson.time_end +
    " ";
  if (props.lesson.free === true) {
    panelTitle += "Свободно";
    return (
      <CollapsablePanel title={`${panelTitle}`} type={"free"}>
        <div>
          <form onSubmit={submitHandler}>
            <textarea
              value={comment}
              placeholder="Комментарий"
              onChange={(event) => setComment(event.target.value)}
              rows="3"
            ></textarea>
            <button type="submit">Забронировать</button>
          </form>
        </div>
      </CollapsablePanel>
    );
  } else if (props.lesson.reserved === true) {
    const { reserver, comment } = props.lesson;
    panelTitle += reserver;
    return (
      <CollapsablePanel title={`${panelTitle}`} type={"reserved_lesson"}>
        <span>Комменатрий: {comment}</span>
      </CollapsablePanel>
    );
  } else {
    const { lector, subject, groups, type } = props.lesson;
    panelTitle += lector;
    return (
      <CollapsablePanel title={`${panelTitle}`} type={"schedule_lesson"}>
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
