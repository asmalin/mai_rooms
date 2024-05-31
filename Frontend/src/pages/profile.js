import AdminNav from "../components/AdminNav";
import "../styles/Profile.css";
import { useEffect, useState } from "react";
import { updateUser, ChangePassword } from "../services/userService.js";

export default function Profile({ user }) {
  const [isEditMode, setIsEditMode] = useState(false);
  const [email, setEmail] = useState("");
  const [tgUsername, setTgUsername] = useState("");

  const [oldPassword, setOldPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [message, setMessage] = useState("");
  const [messageType, setMessageType] = useState("");

  useEffect(() => {
    if (user) {
      setEmail(user.email);
      setTgUsername(
        "@" + (user.tgUsername !== "" ? user.tgUsername : "username")
      );
    }
  }, [user]);

  const tgTagChange = (e) => {
    const newTgUsername = e.target.value;
    if (tgUsername[0] !== "@") {
      setTgUsername("@" + newTgUsername);
    } else {
      setTgUsername(newTgUsername);
    }
  };

  const cancelEditing = () => {
    if (user) {
      setEmail(user.email);
    }
    setIsEditMode(!isEditMode);
  };

  const saveChanges = () => {
    const tgName = tgUsername.substring(1);
    updateUser({ id: user.id, email: email, tgUsername: tgName });
    setIsEditMode(!isEditMode);
  };

  const handleChangePassword = () => {
    if (newPassword !== confirmPassword) {
      setMessage("Новый пароль и подтверждение пароля не совпадают");
      setMessageType("error");
      return;
    }

    ChangePassword(oldPassword, newPassword).then((response) => {
      if (response.ok) {
        setMessage("Пароль успешно изменен");
        setMessageType("success");
        setOldPassword("");
        setNewPassword("");
        setConfirmPassword("");
      } else {
        setMessage("Неправильный старый пароль");
        setMessageType("error");
      }
    });
  };

  return (
    <div className="profile">
      {user?.role === "admin" && <AdminNav />}
      <div className="profile_info wrapper">
        <h1>Профиль</h1>
        <ul className="profile_ul">
          <li>
            <span>Имя:</span>
            <input
              type="text"
              className="profile_field"
              defaultValue={user?.fullname}
              disabled
            />
          </li>
          <li>
            <span>email:</span>
            <input
              type="email"
              className="profile_field"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              disabled={!isEditMode}
              required
            />
          </li>
          <li>
            <span>tg username:</span>
            <input
              type="text"
              className="profile_field"
              value={tgUsername}
              onChange={(e) => tgTagChange(e)}
              disabled={!isEditMode}
            />
          </li>
        </ul>
        <div>
          {isEditMode ? (
            <>
              <button className="btn btn-success" onClick={saveChanges}>
                Сохранить изменения
              </button>
              <button className="btn btn-danger ms-3" onClick={cancelEditing}>
                Отменить
              </button>
            </>
          ) : (
            <button
              className="btn btn-primary"
              onClick={() => setIsEditMode(!isEditMode)}
            >
              Редактировать
            </button>
          )}
        </div>
        <h3 className="mt-4">Сменить пароль</h3>
        <input
          type="password"
          className="profile_field"
          placeholder="Старый пароль"
          value={oldPassword}
          onChange={(e) => setOldPassword(e.target.value)}
          autoComplete="new-password"
        />
        <input
          type="password"
          className="profile_field"
          placeholder="Новый пароль"
          value={newPassword}
          onChange={(e) => setNewPassword(e.target.value)}
          autoComplete="new-password"
        />
        <input
          type="password"
          className="profile_field"
          placeholder="Новый пароль еще раз"
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
          autoComplete="new-password"
        />
        <button className="btn btn-primary" onClick={handleChangePassword}>
          Сменить
        </button>
        {message && (
          <p className={messageType === "error" ? "errorMsg" : "successMsg"}>
            {message}
          </p>
        )}
      </div>
    </div>
  );
}
