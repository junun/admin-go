import styles from './index.css';
import {Spin} from "antd";


export default function() {
  return (
    <div className={styles.normal}>
      <Spin />
    </div>
  );
}
