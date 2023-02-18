import java.util.Scanner;
import java.util.Queue;
import java.util.LinkedList;

class Pair {
    int x, y;

    Pair(int x, int y) {
        this.x = x;
        this.y = y;
    }
}

public class Main {
    static int N, M;
    static char[][] maze;
    static boolean[][] visited;
    static int[] dx = { -1, 0, 1, 0 };
    static int[] dy = { 0, 1, 0, -1 };

    public static void main(String[] args) {
        Scanner sc = new Scanner(System.in);
        N = sc.nextInt();
        M = sc.nextInt();
        sc.nextLine();

        maze = new char[N][M];
        visited = new boolean[N][M];
        Pair start = null, end = null;

        for (int i = 0; i < N; i++) {
            String line = sc.nextLine();
            for (int j = 0; j < M; j++) {
                maze[i][j] = line.charAt(j);
                if (maze[i][j] == 'P') {
                    end = new Pair(i, j);
                } else if (maze[i][j] == 'B') {
                    start = new Pair(i, j);
                }
            }
        }

        int result = bfs(start, end);
        System.out.println(result);
    }

    static int bfs(Pair start, Pair end) {
        Queue<Pair> q = new LinkedList<>();
        q.offer(start);
        visited[start.x][start.y] = true;

        int step = 0;
        while (!q.isEmpty()) {
            int qSize = q.size();
            for (int i = 0; i < qSize; i++) {
                Pair curr = q.poll();
                if (curr.x == end.x && curr.y == end.y) {
                    return step;
                }

                for (int j = 0; j < 4; j++) {
                    int nx = curr.x + dx[j];
                    int ny = curr.y + dy[j];
                    if (nx >= 0 && nx < N && ny >= 0 && ny < M && !visited[nx][ny] && maze[nx][ny] != '#') {
                        q.offer(new Pair(nx, ny));
                        visited[nx][ny] = true;
                    }
                }
            }
            step++;
        }
        return -1;
    }
}