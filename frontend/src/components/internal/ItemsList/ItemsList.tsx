import React from "react";
import Item from "./Item.tsx";

const ItemsList: React.FC = () => {
    return (
        <div className="w-full flex flex-col">
            <Item />
            <Item />
            <Item />
        </div>
    );
}

export default ItemsList;