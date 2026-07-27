import { Card, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";

const items = [
  {
    persona: "바이브코더",
    title: "agent가 만든 코드, 내 이해 속도보다 빠르게 쌓인다",
    description:
      "커밋은 쌓이는데 무슨 기능이 언제 어디에 들어갔는지 따라가기 벅차요.",
  },
  {
    persona: "PM / PO",
    title: "무엇이 왜 바뀌었는지 알려면 개발자를 붙잡고 물어야 한다",
    description:
      "git에는 접근할 수 없고, 프로덕트 현황을 아는 유일한 통로가 사람입니다.",
  },
  {
    persona: "비개발자 창업자",
    title: "내 프로덕트인데 뭘 하는지 모른다",
    description: "코드는 내 자산인데, 그 안에서 무슨 일이 일어나는지는 깜깜해요.",
  },
];

export function Problem() {
  return (
    <section className="px-6 py-16">
      <div className="mx-auto max-w-3xl">
        <h2 className="text-center text-2xl font-semibold tracking-tight">
          기술부채는 익숙하시죠. 이건 새로운 부채입니다
        </h2>
        <p className="mt-3 text-center text-sm leading-relaxed text-muted-foreground">
          레거시 방식에선 코드 짜는 속도가 병목이었습니다. 지금은 AI가 그
          속도를 뛰어넘으면서, 만든 사람조차 못 따라가는 인지부채가 새로
          쌓이고 있습니다.
        </p>
        <div className="mt-10 grid gap-4 sm:grid-cols-3">
          {items.map((item) => (
            <Card key={item.persona}>
              <CardHeader>
                <p className="text-xs font-medium tracking-tight text-muted-foreground">
                  {item.persona}
                </p>
                <CardTitle className="mt-1 text-base leading-snug">
                  {item.title}
                </CardTitle>
                <CardDescription className="mt-2">
                  {item.description}
                </CardDescription>
              </CardHeader>
            </Card>
          ))}
        </div>
      </div>
    </section>
  );
}
