import React, { FunctionComponent } from 'react';
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"

import {MessageCirclePlus} from "lucide-react";
import {Input} from "@/components/ui/input.tsx";
import {Label} from "@/components/ui/label.tsx";

import dayjs from "dayjs"
import GradeList from "@/components/internal/GradeList.tsx";

type TProps = {};

const timerFormat = 'HH:mm';

const CreateInterviewDialog: FunctionComponent<TProps> = (props) => {

  return (
      <Dialog>
        <DialogTrigger asChild>
          <Button><MessageCirclePlus /> Новое интервью</Button>
        </DialogTrigger>
        <DialogContent className="sm:max-w-[425px]">

          <DialogHeader>
            <DialogTitle>Создание интервью</DialogTitle>
          </DialogHeader>
          <div>
              <Label className={"pl-1"}>Название интервью</Label>
              <Input className={"mt-1"} name={"title"} placeholder={"Собеседование в Google"}></Input>

              <div className={"flex items-baseline"}>
                  <div className={"p-2"}>
                        <Label className={"pl-1"}>Тайминг</Label>
                  </div>
                  <div className={"w-24"}>
                    <Input className={"mt-4"} name={"date"} type={"time"}></Input>
                  </div>
              </div>
              <div className={"mt-3"}>
                  <GradeList/>
              </div>
          </div>
          <DialogFooter>
            <Button type="submit">Save changes</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
  );
};

export default CreateInterviewDialog;
