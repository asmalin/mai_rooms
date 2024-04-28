import "../styles/Profile.css";

export default function Profile({ user }) {
  return (
    <>
      <h1>Профиль</h1>
      <ul className="profile_ul">
        <li>
          <span>Имя:</span>
          <input
            type="text"
            className="profile_field"
            defaultValue={user ? user.fullname : ""}
            disabled
          />
        </li>
        <li>
          <span>email:</span>

          <input
            type="text"
            className="profile_field"
            defaultValue={user ? user.email : ""}
            disabled
          />
        </li>
      </ul>
    </>
  );
}
