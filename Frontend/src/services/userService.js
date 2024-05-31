export async function getAllUsers() {
  const token = localStorage.getItem("token");
  if (!token) {
    console.log("Токен не найден");
    return;
  }
  try {
    const response = await fetch("http://localhost:5001/api/users", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${JSON.parse(token)}`,
      },
    });

    if (!response.ok) {
      return "error";
    }

    const result = await response.json();
    return result;
  } catch (error) {
    return error;
  }
}

export async function createUser(userData) {
  const token = localStorage.getItem("token");
  if (!token) {
    console.log("Токен не найден");
    return;
  }
  try {
    const response = await fetch("http://localhost:5001/api/users/create", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${JSON.parse(token)}`,
      },
      body: JSON.stringify(userData),
    });

    return response;
  } catch (error) {
    return error;
  }
}

export async function deleteUser(userId) {
  const token = localStorage.getItem("token");
  if (!token) {
    console.log("Токен не найден");
    return;
  }
  try {
    const response = await fetch(
      "http://localhost:5001/api/users/delete/" + userId,
      {
        method: "DELETE",
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

export async function updateUser(userData) {
  const token = localStorage.getItem("token");
  if (!token) {
    console.log("Токен не найден");
    return;
  }
  try {
    const response = await fetch(
      "http://localhost:5001/api/users/update/" + userData.id,
      {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${JSON.parse(token)}`,
        },
        body: JSON.stringify(userData),
      }
    );

    return response;
  } catch (error) {
    return error;
  }
}

export async function ChangePassword(oldPass, newPass) {
  const token = localStorage.getItem("token");
  if (!token) {
    console.log("Токен не найден");
    return;
  }
  try {
    const pass = {
      oldPassword: oldPass,
      newPassword: newPass,
    };
    const response = await fetch(
      "http://localhost:5001/api/users/update/password",
      {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${JSON.parse(token)}`,
        },
        body: JSON.stringify(pass),
      }
    );

    return response;
  } catch (error) {
    return error;
  }
}
