import React from "react";
import {Separator} from "@/components/ui/separator.tsx";
import {Button} from "@/components/ui/button.tsx";
import { MessageCirclePlus } from "lucide-react";
import {ItemsList} from "@/components/internal/ItemsList";


const List: React.FC = () => {
    return (
        <div className={"flex flex-col w-full"}>
            <div className={"w-full flex p-5 justify-around"}>
                <Button><MessageCirclePlus /> Новое интервью</Button>
            </div>
            <Separator/>
            <div>
                <ItemsList />
            </div>
        </div>
    );
}

export default List;