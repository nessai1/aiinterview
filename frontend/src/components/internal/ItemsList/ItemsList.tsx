import React from "react";
import Item from "./Item.tsx";

const ItemsList: React.FC = () => {
    return (
        <div className="w-full flex flex-col">
            <Item complete={true} />
            <Item complete={false} />
        </div>
    );
}

export default ItemsList;