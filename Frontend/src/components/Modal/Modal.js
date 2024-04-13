import React from "react";
import "./Modal.css";

const Modal = ({ active, setActive, children }) => {
  return (
    <div
      className={active ? "modalz active" : "modalz"}
      onMouseDown={() => setActive(false)}
      onClick={(e) => e.stopPropagation()}
    >
      <div
        className={active ? "modalz__content active" : "modalz__content"}
        onMouseDown={(e) => e.stopPropagation()}
      >
        {children}
      </div>
    </div>
  );
};

export default Modal;
