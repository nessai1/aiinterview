import React from "react";
import {Separator} from "@/components/ui/separator.tsx";
import {ItemsList} from "@/components/internal/ItemsList";
import CreateInterviewDialog from "@/components/internal/CreateInterviewDialog.tsx";


const List: React.FC = () => {
    return (
        <div className={"flex flex-col w-full"}>
            <div className={"w-full flex p-5 justify-around"}>
                <CreateInterviewDialog />
            </div>
            <Separator/>
            <div>
                <ItemsList />
            </div>
        </div>
    );
}

export default List;