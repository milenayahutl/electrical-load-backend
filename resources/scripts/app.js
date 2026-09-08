document.addEventListener("DOMContentLoaded", function () {

    /*
        =====================================
        1. АВТОМАТИЧЕСКОЕ ПРОИГРЫВАНИЕ ВИДЕО
        =====================================

        Когда видео находится в видимой области ленты,
        оно автоматически запускается.

        Остальные видео останавливаются.
    */

    const feedItems = document.querySelectorAll(".feed-item");
    const feedVideos = document.querySelectorAll(".feed-video");

    if (feedItems.length > 0) {

        function pauseAllVideos() {
            feedVideos.forEach(function (video) {
                video.pause();
            });
        }

        const observer = new IntersectionObserver(
            function (entries) {

                entries.forEach(function (entry) {

                    const video = entry.target.querySelector(".feed-video");

                    if (!video) {
                        return;
                    }

                    if (entry.isIntersecting) {

                        pauseAllVideos();

                        video.play().catch(function () {
                            /*
                                Некоторые браузеры запрещают autoplay,
                                если у видео включён звук.

                                У нас в HTML стоит muted,
                                поэтому обычно всё работает.
                            */
                        });

                    } else {

                        video.pause();

                    }
                });
            },
            {
                /*
                    Видео считается активным,
                    если видно хотя бы 65% карточки.
                */
                threshold: 0.65
            }
        );

        feedItems.forEach(function (item) {
            observer.observe(item);
        });
    }


    /*
        =====================================
        2. КНОПКА "СЛЕДУЮЩИЙ"
        =====================================

        Находим следующую карточку в ленте
        и плавно прокручиваем до неё.

        Если текущая карточка последняя,
        возвращаемся к первой.
    */

    const nextButtons = document.querySelectorAll(".next-button");

    nextButtons.forEach(function (button) {

        button.addEventListener("click", function () {

            const currentItem = button.closest(".feed-item");

            if (!currentItem) {
                return;
            }

            let nextItem = currentItem.nextElementSibling;

            /*
                Если следующего элемента нет,
                значит мы дошли до конца ленты.
            */
            if (!nextItem) {

                const feedList = currentItem.parentElement;

                nextItem = feedList.firstElementChild;
            }

            if (nextItem) {

                nextItem.scrollIntoView({
                    behavior: "smooth",
                    block: "start"
                });
            }
        });
    });


    /*
        =====================================
        3. КНОПКА "БОЛЬШЕ / МЕНЬШЕ"
        =====================================

        По умолчанию описание ограничено
        несколькими строками через CSS.

        После нажатия показываем текст полностью.
    */

    const descriptionButtons =
        document.querySelectorAll(".description-toggle");

    descriptionButtons.forEach(function (button) {

        button.addEventListener("click", function () {

            const infoBlock = button.closest(".feed-info");

            if (!infoBlock) {
                return;
            }

            const description =
                infoBlock.querySelector(".feed-description");

            if (!description) {
                return;
            }

            description.classList.toggle("expanded");

            if (description.classList.contains("expanded")) {

                button.textContent = "меньше";

            } else {

                button.textContent = "больше";
            }
        });
    });


    /*
        =====================================
        4. ВЫБОР ФОТО И ВИДЕО В ФОРМЕ
        =====================================

        После выбора файла вместо текста
        "+ добавить фото" / "+ добавить видео"
        показываем имя выбранного файла.
    */

    const fileInputs = document.querySelectorAll(".file-input");

    fileInputs.forEach(function (input) {

        input.addEventListener("change", function () {

            const uploadBlock = input.closest(".upload-block");

            if (!uploadBlock) {
                return;
            }

            const placeholder =
                uploadBlock.querySelector(".upload-placeholder");

            if (!placeholder) {
                return;
            }

            if (input.files.length > 0) {

                placeholder.textContent = input.files[0].name;

            } else {

                /*
                    Если пользователь отменил выбор файла,
                    возвращаем исходную надпись.
                */

                if (input.accept.includes("image")) {
                    placeholder.textContent = "+ добавить фото";
                }

                if (input.accept.includes("video")) {
                    placeholder.textContent = "+ добавить видео";
                }
            }
        });
    });


    /*
        =====================================
        5. РАСЧЁТ НАГРУЗКИ НА СЕТЬ 220 В
        =====================================

        Пользователь вводит мощность в кВт.

        Формула:

            I = P / U

        Но мощность введена в кВт,
        поэтому сначала переводим её в Вт:

            I = P(кВт) * 1000 / 220

        Например:

            2.0 кВт * 1000 / 220 ≈ 9.1 А
    */

    const powerInput = document.querySelector("#power");
    const loadInput = document.querySelector("#load");

    if (powerInput && loadInput) {

        function calculateLoad() {

            const powerKW = Number(powerInput.value);

            /*
                Если поле пустое,
                не показываем никакого результата.
            */

            if (
                powerInput.value === "" ||
                Number.isNaN(powerKW) ||
                powerKW < 0
            ) {

                loadInput.value = "";
                return;
            }

            const loadA = powerKW * 1000 / 220;

            /*
                Оставляем один знак после запятой.
                Например 9.0909... → 9.1
            */

            loadInput.value = loadA.toFixed(1);
        }

        /*
            Пересчитываем нагрузку каждый раз,
            когда пользователь меняет мощность.
        */

        powerInput.addEventListener("input", calculateLoad);

        /*
            Если в поле мощности уже есть значение
            при загрузке страницы,
            сразу рассчитываем нагрузку.
        */

        calculateLoad();
    }

});