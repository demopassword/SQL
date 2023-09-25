## MySQL USING VS ON 차이

#### USING
- 두 테이블간 필드이름을 같은 경우에 사용한다.
    ```sql
    SELECT * FROM `Green` INNER JOIN `Tea` USING(plant) WHERE ~
    ```

#### ON
- 조인시에 컬럼 이름이 다를 경우 ON을 사용한다.

    - 이름이 다를 경우
    ```sql
    SELECT * FROM `Green` INNER JOIN `Tea` ON Green.plant = Tea.name WHERE ~
    ```
    - 이름이 같을 경우
    ```
    SELECT * FROM `Green` INNER JOIN `Tea` ON Green.name = Tea.name WHERE ~
    ```