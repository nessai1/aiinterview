import React from "react";
import {Separator} from "@/components/ui/separator.tsx";
import {Button} from "@/components/ui/button.tsx";


const List: React.FC = () => {
    return (
        <div className={"flex flex-col w-full"}>
            <div className={"w-full flex p-5 justify-around"}>
                <Button>Новое интервью</Button>
            </div>
            <Separator/>
            <div>
                List...
            </div>
        </div>
    );
}

export default List;