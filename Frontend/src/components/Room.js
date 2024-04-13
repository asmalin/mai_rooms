import React, { useEffect, useState } from "react";
import Lesson from "./Lesson.js";
import "../styles/Room.css";

export default function Room({ roomId, date }) {
  const [lessons, setLessons] = useState([]);
  const [roomName, setRoomName] = useState("");

  useEffect(() => {
    let apiUrl = `http://localhost:8080/api/room/${roomId}`;
    fetch(apiUrl)
      .then((response) => response.json())
      .then((data) => setRoomName(data));

    apiUrl = `http://localhost:8080/api/lessons?date=${date}&room=${roomId}`;
    fetch(apiUrl)
      .then((response) => response.json())
      .then((data) => setLessons(data));
  }, [date, roomId]);

  return (
    <>
      <span className="roomName">{roomName}</span>
      <div className="lessons">
        {lessons.map((lesson, index) => {
          return (
            <div key={index} className="lesson">
              <Lesson
                lesson={lesson}
                lessonNumber={index + 1}
                date={date}
                roomId={roomId}
              />
            </div>
          );
        })}
      </div>
    </>
  );
}
