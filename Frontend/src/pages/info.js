import React from "react";

function Info() {
  return (
    <>
      <div>
        Здесь можно посмотреть расписание занятий на аудиторию, а также
        забронировать свободную.
      </div>
      <br />
      <div>
        А так же есть{" "}
        <a
          href="https://t.me/roomsReservation_bot"
          target="_blank"
          rel="noreferrer"
        >
          telegram
        </a>{" "}
        версия приложения.
      </div>
    </>
  );
}

export default Info;
