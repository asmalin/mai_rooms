export async function reserveRoom(lessonForReservation) {
  const token = localStorage.getItem("token");
  if (!token) {
    console.log("Токен не найден");
    return;
  }
  try {
    const response = await fetch("http://localhost:8080/api/reserve", {
      method: "POST",
      body: JSON.stringify(lessonForReservation),
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${JSON.parse(token)}`,
      },
    });

    if (!response.ok) {
      return "error";
    }

    const data = await response.json();
    return data;
  } catch (error) {
    return error;
  }
}

export async function cancelReserve(lessonForCancelReservation) {
  const token = localStorage.getItem("token");
  if (!token) {
    console.log("Токен не найден");
    return;
  }
  try {
    const response = await fetch(
      "http://localhost:8080/api/cancelReservation",
      {
        method: "POST",
        body: JSON.stringify(lessonForCancelReservation),
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${JSON.parse(token)}`,
        },
      }
    );

    if (!response.ok) {
      return "error";
    }

    const result = await response.json();
    return result;
  } catch (error) {
    return error;
  }
}

export async function CreateQRCodes(selectedRooms) {
  const token = localStorage.getItem("token");
  if (!token) {
    console.log("Токен не найден");
    return;
  }

  let roomIds = "roomId=";
  roomIds += selectedRooms.join("&roomId=");

  try {
    const response = await fetch(
      "http://localhost:8080/api/qrcodes?" + roomIds,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${JSON.parse(token)}`,
        },
      }
    );

    if (!response.ok) {
      throw new Error("Ошибка при загрузке архива");
    }

    const archiveBlob = await response.blob();
    const link = document.createElement("a");
    link.href = URL.createObjectURL(archiveBlob);
    link.download = "qrcodes.zip";
    link.click();
    URL.revokeObjectURL(archiveBlob);
    return;
  } catch (error) {
    return error;
  }
}
