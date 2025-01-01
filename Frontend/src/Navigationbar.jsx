import USnameandpp from "./USnameandpp";
import './Navbar.css'
import Navbuttons from "./Navbuttons";
export default function Navigationbar(){
    return(
        <div className="nav-bar">
<USnameandpp/>
<hr className="three-d-rule" />
<Navbuttons/>


</div>
    )
}