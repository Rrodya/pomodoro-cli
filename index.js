import inquirer from 'inquirer';
let workTime = 25; 
let breakTime = 5;
let timer = null;
let isPaused = false;
let isWorking = true;
let remainingTime = 0;

const formatTime = (seconds) => {
  const minutes = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${String(minutes).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
};

const handleUserInput = () => {
  process.stdin.setRawMode(true);
  process.stdin.resume();
  process.stdin.setEncoding('utf8');

  process.stdin.on('data', (key) => {
    if (key === 'p') { 
      if (!isPaused && timer) {
        isPaused = true;
        console.log('Таймер поставлен на паузу. Нажмите "r" для продолжения.');
      } else {
        console.log('Таймер уже на паузе.');
      }
    } else if (key === 'r') {
      if (isPaused) {
        isPaused = false;
        console.log('Таймер продолжен.');
      } else {
        console.log('Таймер уже работает.');
      }
    } else if (key === 's') { 
      if (timer) {
        clearInterval(timer);
        timer = null;
        isPaused = false;
        console.log('Таймер остановлен.');
        menu(); 
      } else {
        console.log('Таймер не запущен.');
      }
    } else if (key === '\u0003') {
      console.log('До свидания!');
      process.exit();
    }
  });
};

const startTimer = (duration, callback) => {
  remainingTime = duration;
  console.log(`Таймер запущен: ${formatTime(remainingTime)} осталось.`);
  timer = setInterval(() => {
    if (!isPaused) {
      remainingTime--;
      console.clear();
      console.log(`${formatTime(remainingTime)} осталось. Нажмите "p" для паузы, "r" для продолжения, "s" для остановки.`);
      if (remainingTime <= 0) {
        clearInterval(timer);
        timer = null;
        console.log(`Время ${isWorking ? 'работы' : 'перерыва'} закончилось!`);
        callback();
      }
    }
  }, 1000);
};

const menu = async () => {
  const { choice } = await inquirer.prompt([
    {
      type: 'list',
      name: 'choice',
      message: 'Выберите действие:',
      choices: [
        { name: `Установить время работы (текущее: ${workTime} мин)`, value: 'setWorkTime' },
        { name: `Установить время перерыва (текущее: ${breakTime} мин)`, value: 'setBreakTime' },
        { name: `Запустить таймер`, value: 'startTimer' },
        { name: `Выход`, value: 'exit' },
      ],
    },
  ]);

  switch (choice) {
    case 'setWorkTime':
      const { work } = await inquirer.prompt([
        { type: 'number', name: 'work', message: 'Введите время работы (в минутах):' },
      ]);
      workTime = work;
      console.log(`Время работы установлено: ${workTime} мин.`);
      break;

    case 'setBreakTime':
      const { break: br } = await inquirer.prompt([
        { type: 'number', name: 'break', message: 'Введите время перерыва (в минутах):' },
      ]);
      breakTime = br;
      console.log(`Время перерыва установлено: ${breakTime} мин.`);
      break;

    case 'startTimer':
      if (timer) {
        console.log('Таймер уже запущен!');
      } else {
        isWorking = true;
        handleUserInput();
        startTimer(workTime * 60, () => {
          isWorking = false;
          startTimer(breakTime * 60, () => {
            console.log('Цикл завершен!');
            menu();
          });
        });
      }
      return;

    case 'exit':
      console.log('До свидания!');
      process.exit(0);

    default:
      console.log('Неверный выбор.');
  }

  menu();
};

menu();
