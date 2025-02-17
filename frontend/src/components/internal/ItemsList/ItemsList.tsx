import React, {useEffect, useState} from "react";
import Item from "./Item.tsx";
import {Interview} from "@/lib/interview/interview.ts";
import {Skeleton} from "@/components/ui/skeleton.tsx";
import {useToast} from "@/hooks/use-toast.ts";
import {ToastAction} from "@/components/ui/toast.tsx";
import {AxiosError} from "axios";
import {Link} from "react-router-dom";

const ItemsList: React.FC = () => {
    const [items, setItems] = useState<Interview[]>()
    const { toast } = useToast();


    const retry = () => {
        network.getInterviewList().then((interviews: Interview[]) => {
            setItems(interviews);
        }).catch((err: AxiosError) => {
            toast({
                title: 'Упс! Список интервью не загрузился ;(',
                description: `Ошибка сети: [${err.code}] ${err.message}`,
                action: (
                    <ToastAction altText="Goto schedule to undo" onClick={retry}>Повторить</ToastAction>
                ),
                variant: "destructive",
            });
        });
    };

    useEffect(() => {
        retry();
    }, []);

    return (
        <div className="w-full flex flex-col">
            {items === undefined ?
                <div className="flex flex-col space-y-2 pt-2">
                    <Skeleton className="w-full h-[100px] rounded-lg bg-zinc-900" />
                    <Skeleton className="w-full h-[100px] rounded-lg bg-zinc-900" />
                    <Skeleton className="w-full h-[100px] rounded-lg bg-zinc-900" />
                </div>
                :
                <div>
                    { items.length > 0
                        ? items.map((val) => <Link to={"/interview/" + val.uuid} key={val.uuid} style={{ textDecoration: "none", color: "inherit" }}><Item interview={val} /></Link>)
                        : <h3 className={"mt-10"}>Тут пока ничего нет :)</h3>
                    }
                </div>}
        </div>
    );
}

export default ItemsList;