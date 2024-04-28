import React, { useEffect, useState } from "react";
import Lesson from "./Lesson.js";
import "../styles/Room.css";
import {
  getLessonNumber,
  lessonTimeStart,
  lessonTimeEnd,
} from "../services/lessonsService.js";

export default function Room({ roomId, date, showOnlyFree, user }) {
  const [lessons, setLessons] = useState([]);
  const [roomName, setRoomName] = useState("");

  const [scheduleLessons, setScheduleLessons] = useState([]);
  const [reservedLessons, setReservedLessons] = useState([]);

  useEffect(() => {
    fetch(`http://localhost:8080/api/room/${roomId}`)
      .then((response) => response.json())
      .then((data) => setRoomName(data));

    fetch(`http://localhost:8080/api/schedule?date=${date}&room=${roomId}`)
      .then((response) => response.json())
      .then((lessons) => setScheduleLessons(getLessonNumber(lessons)));

    fetchReservedRooms();
  }, [date, roomId]);

  function fetchReservedRooms() {
    fetch(
      `http://localhost:8080/api/reserved_lesssons?date=${date}&room=${roomId}`
    )
      .then((response) => response.json())
      .then((lessons) => setReservedLessons(getLessonNumber(lessons)));
  }

  return (
    <>
      <span className="roomName">{roomName}</span>
      <div className="lessons">
        {Object.entries(lessonTimeStart).map(([timestart, lessonNumber]) => {
          let lesson =
            scheduleLessons[lessonNumber] ||
            reservedLessons[lessonNumber] ||
            null;
          if (showOnlyFree && lesson !== null) {
            return null;
          }

          return (
            <div key={lessonNumber} className="lesson">
              <Lesson
                lessonNumber={lessonNumber}
                roomId={roomId}
                date={date}
                timeStart={timestart}
                timeEnd={lessonTimeEnd[lessonNumber]}
                lesson={lesson}
                fetchReservedRooms={fetchReservedRooms}
                user={user}
              />
            </div>
          );
        })}
      </div>
    </>
  );
}
