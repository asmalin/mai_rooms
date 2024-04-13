import React, { useState } from "react";
import "../styles/CollapsablePanel.css";

const CollapsiblePanel = ({ title, type, children }) => {
  const [isCollapsed, setIsCollapsed] = useState(true);

  const togglePanel = () => {
    setIsCollapsed(!isCollapsed);
  };

  return (
    <div className={type}>
      <div className="panelHeader" onClick={togglePanel}>
        <span>{title}</span>
        <span>{isCollapsed ? "▼" : "▲"}</span>
      </div>
      {!isCollapsed && (
        <div className={`panelBottom ${type}_bottom`}>{children}</div>
      )}
    </div>
  );
};

export default CollapsiblePanel;
