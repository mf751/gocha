import { MdPerson, MdPersonOutline, MdSearch } from "react-icons/md";
import { IoChatbubbleSharp, IoChatbubbleOutline } from "react-icons/io5";
import "./styles.css";
import { Link } from "react-router-dom";

export default function Nav() {
  const current = window.location.pathname.slice(1);
  return (
    <nav>
      <Link to="/">
        {current == "" ? (
          <IoChatbubbleSharp className="icon" />
        ) : (
          <IoChatbubbleOutline className="icon" />
        )}
      </Link>
      <Link to="/search">
        <MdSearch className="icon" />
      </Link>
      <Link to="/profile">
        {current === "profile" ? (
          <MdPerson className="icon" />
        ) : (
          <MdPersonOutline className="icon" />
        )}
      </Link>
    </nav>
  );
}
