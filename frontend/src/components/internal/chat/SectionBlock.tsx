import React from "react";
import {Section} from "@/lib/interview/interview.ts";
import './SectionBlock.css';
import QuestionBlock from "@/components/internal/chat/QuestionBlock.tsx";


type TProps = {
    section: Section
}

const testQuestion = {
    uuid: '0202020202',
    text: `
    <div class="markdown-content"><p>Если бабушка была дедушкой, то как бы выглядел код на go?</p>

<p></p><pre class="chroma"><code><span class="line"><span class="cl">
</span></span><span class="line"><span class="cl"><span class="kd">struct</span> <span class="nx">Ded</span> <span class="p">{</span>
</span></span><span class="line"><span class="cl">    <span class="kt">int</span> <span class="nx">HuiSize</span>
</span></span><span class="line"><span class="cl"><span class="p">}</span>
</span></span><span class="line"><span class="cl">
</span></span><span class="line"><span class="cl"><span class="kd">func</span> <span class="nf">main</span><span class="p">(</span><span class="p">)</span> <span class="p">{</span>
</span></span><span class="line"><span class="cl">    <span class="nx">babushka</span> <span class="o">:=</span> <span class="nx">Ded</span><span class="p">{</span><span class="p">}</span>
</span></span><span class="line"><span class="cl"><span class="p">}</span></span></span></code></pre><p></p>

<p>Назови</p>

<ul>
<li>принципы солид</li>
<li>принципы KISS</li>
</ul>
</div>
`,
    answer: '',
    feedback: 'ну ты ваще долбоеб',
};

const SectionBlock: React.FC<TProps> = (props: TProps) => {

    return (
        <div className="sectionBlock">
            <div className="sectionHeader">
                <div className="circleNumber">1</div>
                <div className="sectionTitle">Алгоритмы</div>
                <div className="sectionGrade middle">middle</div>
            </div>
            <div className="sectionQuestions">
                <QuestionBlock question={testQuestion} showFeedback={false}/>
            </div>
        </div>
    );
}

export default SectionBlock;