import React, { useEffect } from "react";
import { useState } from "react";
import AdminNav from "../components/AdminNav";
import {
  getAllUsers,
  createUser,
  deleteUser,
  updateUser,
} from "../services/userService";

import "../styles/Users.css";

export default function Users({ user }) {
  const [users, setUsers] = useState();
  const [editedIndex, setEditedIndex] = useState(null);
  const [editedUserData, setEditedUserData] = useState({});

  const [newUserData, setNewUserData] = useState({
    username: "",
    fullname: "",
    password: "",
    email: "",
    role: "user",
  });

  const [filter, setFilter] = useState({
    username: "",
    fullname: "",
    role: "",
    email: "",
  });

  const editUserHandler = (index, user) => {
    setEditedUserData(user);
    setEditedIndex(index);
  };

  const cancelEditHandler = () => {
    setEditedUserData({});
    setEditedIndex(null);
  };

  const saveUserHandler = () => {
    updateUser(editedUserData).then((response) => {
      if (response.ok) {
        setEditedUserData({});
        setEditedIndex(null);
        getAllUsers().then((users) => setUsers(users));
      } else {
        console.log(response);
      }
    });
  };

  const handleChange = (e) => {
    setNewUserData({ ...newUserData, [e.target.name]: e.target.value });
  };

  useEffect(() => {
    getAllUsers().then((users) => setUsers(users));
  }, []);

  const filteredUsers = users
    ? users.filter((user) => {
        return (
          user.username.toLowerCase().includes(filter.username.toLowerCase()) &&
          user.fullname.toLowerCase().includes(filter.fullname.toLowerCase()) &&
          user.role.toLowerCase().includes(filter.role.toLowerCase()) &&
          user.email.toLowerCase().includes(filter.email.toLowerCase())
        );
      })
    : [];

  const handleSubmit = (e) => {
    e.preventDefault();

    createUser(newUserData).then((response) => {
      if (response.ok) {
        getAllUsers().then((users) => setUsers(users));
        setNewUserData({
          username: "",
          fullname: "",
          password: "",
          email: "",
          role: "user",
        });
      }
    });
  };

  const deleteUserHandler = (userId) => {
    deleteUser(userId).then(() =>
      getAllUsers().then((users) => setUsers(users))
    );
  };

  return (
    <>
      <AdminNav />
      <div className="users_mgmt">
        <strong>Добавить нового пользователя</strong>
        <form onSubmit={handleSubmit}>
          <input
            placeholder="Имя пользователя"
            name="username"
            className="form-control"
            value={newUserData.username}
            onChange={handleChange}
            autoComplete="new-password"
          />
          <input
            placeholder="Полное имя"
            name="fullname"
            className="form-control"
            value={newUserData.fullname}
            onChange={handleChange}
            autoComplete="new-password"
          />
          <input
            placeholder="Пароль"
            name="password"
            type="password"
            className="form-control"
            value={newUserData.password}
            onChange={handleChange}
            autoComplete="new-password"
          />
          <input
            placeholder="Email"
            name="email"
            className="form-control"
            value={newUserData.email}
            onChange={handleChange}
            autoComplete="new-password"
          />
          <select
            name="role"
            className="form-select"
            value={newUserData.role}
            onChange={handleChange}
            autoComplete="new-password"
          >
            <option value="user">user</option>
            <option value="admin">admin</option>
          </select>
          <button type="submit">Добавить</button>
        </form>
        <h3>Все пользователи</h3>

        <div>
          <strong>Фильтры</strong>
          <form className="filter_form">
            <input
              className="form-control form-control-lg"
              type="text"
              placeholder="username"
              value={filter.username}
              onChange={(e) =>
                setFilter({ ...filter, username: e.target.value })
              }
            />
            <input
              className="form-control form-control-lg"
              type="text"
              placeholder="Имя"
              value={filter.fullname}
              onChange={(e) =>
                setFilter({ ...filter, fullname: e.target.value })
              }
            />
            <input
              className="form-control form-control-lg"
              type="text"
              placeholder="email"
              value={filter.email}
              onChange={(e) => setFilter({ ...filter, email: e.target.value })}
            />
            <select
              className="form-select"
              onChange={(e) =>
                e.target.value !== "Все роли"
                  ? setFilter({ ...filter, role: e.target.value })
                  : setFilter({ ...filter, role: "" })
              }
            >
              <option>Все роли</option>
              <option>admin</option>
              <option>user</option>
            </select>
          </form>
        </div>
        {filteredUsers.length > 0 ? (
          <table className="table table-bordered">
            <thead>
              <tr>
                <th scope="col">#</th>
                <th scope="col">username</th>
                <th scope="col">fullname</th>
                <th scope="col">email</th>
                <th scope="col">role</th>
                <th scope="col"></th>
                <th scope="col"></th>
              </tr>
            </thead>
            <tbody>
              {filteredUsers.map((user, index) => (
                <tr key={index}>
                  <th scope="row">{index + 1}</th>
                  <td>
                    {editedIndex === index ? (
                      <input
                        className="table_cell_editMode"
                        value={editedUserData.username}
                        onChange={(e) =>
                          setEditedUserData({
                            ...editedUserData,
                            username: e.target.value,
                          })
                        }
                      />
                    ) : (
                      user.username
                    )}
                  </td>
                  <td>
                    {editedIndex === index ? (
                      <input
                        className="table_cell_editMode"
                        value={editedUserData.fullname}
                        onChange={(e) =>
                          setEditedUserData({
                            ...editedUserData,
                            fullname: e.target.value,
                          })
                        }
                      />
                    ) : (
                      user.fullname
                    )}
                  </td>
                  <td>
                    {editedIndex === index ? (
                      <input
                        className="table_cell_editMode"
                        value={editedUserData.email}
                        onChange={(e) =>
                          setEditedUserData({
                            ...editedUserData,
                            email: e.target.value,
                          })
                        }
                      />
                    ) : (
                      user.email
                    )}
                  </td>
                  <td>
                    {editedIndex === index ? (
                      <input
                        className="table_cell_editMode"
                        value={editedUserData.role}
                        onChange={(e) =>
                          setEditedUserData({
                            ...editedUserData,
                            role: e.target.value,
                          })
                        }
                      />
                    ) : (
                      user.role
                    )}
                  </td>
                  <td align="center">
                    {editedIndex === index ? (
                      <button
                        className="btn btn-primary"
                        onClick={saveUserHandler}
                      >
                        Сохранить
                      </button>
                    ) : (
                      <button
                        className="btn btn-success"
                        onClick={() => editUserHandler(index, user)}
                      >
                        Редактировать
                      </button>
                    )}
                  </td>
                  <td align="center">
                    {editedIndex === index ? (
                      <button
                        className="btn btn-secondary"
                        onClick={cancelEditHandler}
                      >
                        Отмена
                      </button>
                    ) : (
                      <button
                        className="btn btn-danger"
                        onClick={() => deleteUserHandler(user.id)}
                      >
                        Удалить
                      </button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : (
          <h2>Список пуст!</h2>
        )}
      </div>
    </>
  );
}
